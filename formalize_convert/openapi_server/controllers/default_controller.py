import connexion
from typing import Dict
from typing import Tuple
from typing import Union

from openapi_server.models.formalize_request import FormalizeRequest  # noqa: E501
from openapi_server.models.formalize_result import FormalizeResult  # noqa: E501
from openapi_server.models.fix_request import FixRequest  # noqa: E501
from openapi_server.models.fix_response import FixResponse  # noqa: E501
from openapi_server import util

import main as logic


def healthcheck():  # noqa: E501
    """Healthcheck

     # noqa: E501


    :rtype: Union[None, Tuple[None, int], Tuple[None, int, Dict[str, str]]
    """
    return "OK", 200, {"Content-Type": "text/plain"}


def trs_formalize(body):  # noqa: E501
    """Extract formal definition of trs

     # noqa: E501

    :param formalize_request:
    :type formalize_request: dict | bytes

    :rtype: Union[FormalizeResult, Tuple[FormalizeResult, int], Tuple[FormalizeResult, int, Dict[str, str]]
    """
    if connexion.request.is_json:
        formalize_request = FormalizeRequest.from_dict(connexion.request.get_json())  # noqa: E501
        res = logic.formalize(formalize_request.trs)
        if res is None:
            return "Unexpected error", 500, {"Content-Type": "text/plain"}
        return FormalizeResult(res["formalTrs"], res["error"]), 200, {"Content-Type": "application/json"}

    return "Bad request", 400, {"Content-Type": "text/plain"}


def trs_fix(body):  # noqa: E501
    """Fix extracted formal definition of trs

     # noqa: E501

    :param fix_request:
    :type fix_request: dict | bytes

    :rtype: Union[FixResponse, Tuple[FixResponse, int], Tuple[FixResponse, int, Dict[str, str]]
    """
    if connexion.request.is_json:
        fix_request = FixRequest.from_dict(connexion.request.get_json())  # noqa: E501
        trs, formalTrs, error = fix_request.trs, fix_request.formal_trs, fix_request.error
        res = logic.fix_formalized_trs(trs, formalTrs, error)
        if res is None:
            return "Unexpected error", 500, {"Content-Type": "text/plain"}
        return FixResponse(res["formalTrs"], res["error"]), 200, {"Content-Type": "application/json"}

    return "Bad request", 400, {"Content-Type": "text/plain"}
