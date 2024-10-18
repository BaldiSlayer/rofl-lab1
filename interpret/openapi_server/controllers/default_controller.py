import connexion
from typing import Dict
from typing import Tuple
from typing import Union

from openapi_server.models.interpret_result import InterpretResult  # noqa: E501
from openapi_server.models.trs import Trs  # noqa: E501
from openapi_server import util

import logic.interpret as logic


def interpret(trs: Trs) -> str:
    return logic.interpret(trs_variables=trs.variables, trs_rules=trs.rules, grammar_rules=trs.interpretations)


def trs_interpret(body):  # noqa: E501
    """Check decidability of trs

     # noqa: E501

    :param trs: Trs
    :type trs: dict | bytes

    :rtype: Union[InterpretResult, Tuple[InterpretResult, int], Tuple[InterpretResult, int, Dict[str, str]]
    """
    if connexion.request.is_json:
        trs = Trs.from_dict(connexion.request.get_json())  # noqa: E501
        return InterpretResult(interpret(trs)), 200, {"Content-Type": "application/json"}

    return "Bad request", 400, {"Content-Type": "text/plain"}


def healthcheck():  # noqa: E501
    """Healthcheck

     # noqa: E501


    :rtype: Union[None, Tuple[None, int], Tuple[None, int, Dict[str, str]]
    """
    return "OK", 200, {"Content-Type": "text/plain"}
