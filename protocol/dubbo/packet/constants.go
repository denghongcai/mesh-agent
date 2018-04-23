package packet

const DUBBO_VERSION_KEY = "dubbo"
const PATH_KEY = "path"
const VERSION_KEY = "version"

const DUBBO_VERSION = "2.5.3"
const HEADER_LENGTH = 16
const MAGIC_HIGH = 0xda
const MAGIC_LOW = 0xbb
const FLAG_EVENT = 0x20
const FLAG_TWOWAY = 0x40
const FLAG_REQUEST = 0x80
const HEARTBEAT_EVENT = 1
const SERIALIZATION_MASK = 0x1f

// Response
const RESPONSE_NULL_VALUE = 2
const RESPONSE_VALUE = 1
const RESPONSE_WITH_EXCEPTION = 0
// Response Status
const RESPONSE_OK = 20
const RESPONSE_CLIENT_TIMEOUT = 30
const RESPONSE_SERVER_TIMEOUT = 31
const RESPONSE_BAD_REQUEST = 40
const RESPONSE_BAD_RESPONSE = 50
const RESPONSE_SERVICE_NOT_FOUND = 60
const RESPONSE_SERVICE_ERROR = 70
const RESPONSE_SERVER_ERROR = 80
const RESPONSE_CLIENT_ERROR = 90