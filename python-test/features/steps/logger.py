import logging
from colorlog import ColoredFormatter
from configs import TestConfig

configs = TestConfig.configs()


class Logger:
    _LOG_DEFAULT_FORMAT = "  %(log_color)s%(levelname)-8s%(reset)s | %(asctime)s | %(log_color)s%(message)s%(reset)s"

    def __init__(self):
        log_level_options = \
            {"debug": logging.DEBUG,
             "info": logging.INFO,
             "warning": logging.WARNING,
             "error": logging.ERROR,
             "critical": logging.CRITICAL}
        log_level = log_level_options[configs["log_level"]]
        self._init_logger(log_level)

    def _init_logger(self, log_level):
        logging.root.setLevel(log_level)
        formatter = ColoredFormatter(self._LOG_DEFAULT_FORMAT)
        stream = logging.StreamHandler()
        stream.setLevel(log_level)
        stream.setFormatter(formatter)
        self._logger = logging.getLogger("pythonConfig")
        self._logger.setLevel(log_level)
        if self._logger.hasHandlers():
            self._logger.handlers.clear()
        if configs["stream_logs"] == "true":
            self._logger.addHandler(stream)

    def logger_instance(self):
        return self._logger
