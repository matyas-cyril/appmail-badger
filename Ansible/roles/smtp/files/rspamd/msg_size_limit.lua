-- MTK 30 juillet 2026

local rspamd_logger = require "rspamd_logger"
local rspamd_util = require "rspamd_util"
local lua_util = require("lua_util")

-- Charger la configuration dans local.d/sender_size_limit.conf
local config = rspamd_config:get_all_opt("sender_size_limit")

-- Variables par défaut
local SCRIPT_PATH = (config and config.script_path) or "/usr/local/bin/rspamd_pg_query.py"
local SCRIPT_ARG = (config and config.script_arg) or "-maxsize"

local CACHE_TTL = (config and config.cache_ttl) or 15 -- en secondes
local CACHE_FREQ_CLEAN = 100

local SYMBOL_NAME = (config and config.symbol_name) or "SENDER_SIZE_LIMIT_EXCEEDED"
local SYMBOL_SCORE = (config and config.symbol_score) or 0.0 -- valeur numérique de pénalité attribuée au message 

local MESSAGE_EXCEEDS_LIMIT = (config and config.message_exceeds_limit) or "Message size exceeded size[%d/%d] octets - sender[%s]"
local MESSAGE_QUERY_FAILED = (config and config.message_query_failed) or "Failed to check size limit"
local MESSAGE_NO_SCRIPT = (config and config.message_no_script) or "Invalid script or configuration"

-- Cache en mémoire avec expiration
local cache = {}

-- Pour le traitement du JSon
local ucl = require("ucl")

-- Vider le cache dépassé
local function cleanup_expired_cache()

    local now = rspamd_util.get_time()
    for sender, cached in pairs(cache) do
        if now >= cached.expires then
            cache[sender] = nil
        end
    
    end
end

-- Parcourir le JSon
local function parse_json_response(result)

    if not result or #result == 0 then
        return nil
    end
    
    local parser = ucl.parser()
    local ok, err = parser:parse_string(result)
    
    if not ok then
        -- rspamd_logger:warnx("JSON parse error: %s", err)
        return nil
    end
    
    return parser:get_object()
end

--- 
--- Obtenir la valeur max d'un message pour un Sender
--- @return
---        >= 0 : c'est la valeur qui est en BDD
---        -1 : le Sender n'existe pas
---        -2 : erreur du script
local function get_size_limit(sender_email)

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
        rspamd_logger.errx("Script %s not exist", SCRIPT_PATH)
        return -2  -- Fichier n'existe pas
    end
    file:close()

    -- Requête à la BDD via script externe
    local cmd = string.format("%s %s '%s'",SCRIPT_PATH,SCRIPT_ARG,sender_email:gsub("'", "'\\''"))
    rspamd_logger.debugx("Cmd: %s", cmd)
    
    local handle = io.popen(cmd)
    if not handle then
        rspamd_logger.errx("Failed to access cmd: %s", cmd)
        return -2
    end

    local result = handle:read("*a"):gsub("\n", "")
    local success, exit_type, exit_code = handle:close()

    -- Le script a échoué !!!
    if not success or not result then
        rspamd_logger.errx("Failed cmd: %s", cmd)
        return -2
    end

    result = result:match("^%s*(.-)%s*$") or result
    rspamd_logger.debugx("Result cmd: %s", result)

    -- Traiter la réponse JSon
    if result then

        rspamd_logger.debugx("Cast to JSon")

        local json_data = parse_json_response(result)
        -- On vérifie la structure complète du JSON reçu
        if json_data and json_data.success == true and json_data.maxsize then

            local max_size = tonumber(json_data.maxsize)
            
            if max_size and max_size > 0 then
                cache[sender_email] = {
                    limit = max_size,
                    expires = now + CACHE_TTL
                }
                return max_size
            end
        end
    end

    -- Pas de réponse valide
    return -1
    
end

local function check_sender_size_limit(task)

    local from = task:get_from(1)
    if not from or not from[1] then
        return nil
    end

    local sender_email = from[1]['addr']
    local message_size = task:get_size() -- Obtenir la taille du message

    -- On va vérifier le taille, si disponible
    local max_size = get_size_limit(sender_email)

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
