import logging
from colorlog import ColoredFormatter

configs = PlatformConfig.configs()


class Logger:
    _LOG_DEFAULT_LEVEL = logging.DEBUG
    _LOG_DEFAULT_FORMAT = "  %(log_color)s%(levelname)-8s%(reset)s | %(asctime)s | %(log_color)s%(message)s%(reset)s"

    _config = None
    _logger = None

    def __init__(self):
        self._init_logger()

    def _init_logger(self):
        logging.root.setLevel(self._LOG_DEFAULT_LEVEL)
        formatter = ColoredFormatter(self._LOG_DEFAULT_FORMAT)
        stream = logging.StreamHandler()
        stream.setLevel(self._LOG_DEFAULT_LEVEL)
        stream.setFormatter(formatter)
        self._logger = logging.getLogger("pythonConfig")
        self._logger.setLevel(self._LOG_DEFAULT_LEVEL)
        if self._logger.hasHandlers():
            self._logger.handlers.clear()
        if configs["print_logs"] == "true":
            self._logger.addHandler(stream)

    def logger_instance(self):
        return self._logger
