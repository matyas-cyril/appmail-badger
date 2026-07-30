-- MTK 30 juillet 2026

local rspamd_logger = require "rspamd_logger"
local rspamd_util = require "rspamd_util"

-- Charger la configuration dans local.d/sender_size_limit.conf
local config = rspamd_config:get_all_opt("sender_size_limit")

-- Variables par défaut
local SCRIPT_PATH = (config and config.script_path) or "/usr/local/bin/rspamd_pg_query.py"
local CACHE_TTL = (config and config.cache_ttl) or 15 -- en secondes
local CACHE_FREQ_CLEAN = 100
local SYMBOL_NAME = (config and config.symbol_name) or "SENDER_SIZE_LIMIT_EXCEEDED"
local SYMBOL_SCORE = (config and config.symbol_score) or 0.0 -- valeur numérique de pénalité attribuée au message 
local MESSAGE_EXCEEDS_LIMIT = (config and config.message_exceeds_limit) or "Message size exceeded size[%d/%d] octets - sender[%s]"
local MESSAGE_QUERY_FAILED = (config and config.message_query_failed) or "Failed to check size limit"
local MESSAGE_NO_SCRIPT = (config and config.message_no_script) or "Invalid script or configuration"

-- Cache en mémoire avec expiration
local cache = {}

-- Vider le cache dépassé
local function cleanup_expired_cache()

    local now = rspamd_util.get_time()
    for sender, cached in pairs(cache) do
        if now >= cached.expires then
            cache[sender] = nil
        end
    
    end
end

local function get_sender_limit(sender_email)

    local now = rspamd_util.get_time()

    -- Vérifier le cache et nettoyer si expiré
    if cache[sender_email] then
        local cached = cache[sender_email]
        if now < cached.expires then
            return cached.limit
        else
            cache[sender_email] = nil
        end
    end

    -- Vérifier la présence du script
    local file = io.open(SCRIPT_PATH, "r")
    if not file then
        return -2  -- Fichier n'existe pas
    end
    file:close()

    -- Requête à la BDD via script externe
    local cmd = string.format("%s '%s'",SCRIPT_PATH,sender_email:gsub("'", "'\\''"))
    
    local handle = io.popen(cmd)
    if not handle then
        return -2
    end
    
    local result = handle:read("*a"):gsub("\n", "")
    local success, exit_type, exit_code = handle:close()

    -- Le script a échoué !!!
    if not success or not result then
        return -2
    end

    result = result:match("^%s*(.-)%s*$") or result

    if result and result ~= "-1" then

        local max_size = tonumber(result)
        if max_size then
            -- Mise en cache
            cache[sender_email] = {
                limit = max_size,
                expires = now + CACHE_TTL
            }
            return max_size
        end
    end

end

local function check_sender_size_limit(task)

    local from = task:get_from(1)
    if not from or not from[1] then
        return nil
    end

    local sender_email = from[1]['addr']
    local message_size = task:get_size() -- Obtenir la taille du message

    -- On va vérifier le taille, si disponible
    local max_size = get_sender_limit(sender_email)

    -- Trop Gros
    if max_size >=0 and message_size > max_size then
        -- REJETER le message
        task:set_pre_result({
            action = 'reject',
            score = SYMBOL_SCORE,
            message = string.format(MESSAGE_EXCEEDS_LIMIT, message_size, max_size, sender_email)
        })

        return true
    end

    if max_size == -1 then
        -- REJETER le message
        task:set_pre_result({
            action = 'reject',
            score = SYMBOL_SCORE,
            message = MESSAGE_QUERY_FAILED
        })

        return true
    end

    if max_size == -2 then
        -- REJETER le message
        task:set_pre_result({
            action = 'reject',
            score = SYMBOL_SCORE,
            message = MESSAGE_NO_SCRIPT
        })

        return true
    end

    -- Nettoyage périodique du cache (toutes les 100 requêtes environ)
    if math.random(1, CACHE_FREQ_CLEAN) == 1 then
        cleanup_expired_cache()
    end

    return nil
end

rspamd_config:register_symbol({
    name = SYMBOL_NAME,
    callback = check_sender_size_limit,
    score = SYMBOL_SCORE,
    description = 'Vérification de la limite de taille par expéditeur',
    type = 'callback',
})
