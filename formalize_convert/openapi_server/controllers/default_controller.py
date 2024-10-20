import connexion
from typing import Dict
from typing import Tuple
from typing import Union

from openapi_server.models.formalize_request import FormalizeRequest  # noqa: E501
from openapi_server.models.formalize_result import FormalizeResult  # noqa: E501
from openapi_server import util


def healthcheck():  # noqa: E501
    """Healthcheck

     # noqa: E501


    :rtype: Union[None, Tuple[None, int], Tuple[None, int, Dict[str, str]]
    """
    return "OK", 200, {"Content-Type": "text/plain"}


def trs_formalize(formalize_request):  # noqa: E501
    """Extract formal definition of trs

     # noqa: E501

    :param formalize_request:
    :type formalize_request: dict | bytes

    :rtype: Union[FormalizeResult, Tuple[FormalizeResult, int], Tuple[FormalizeResult, int, Dict[str, str]]
    """
    if connexion.request.is_json:
        formalize_request = FormalizeRequest.from_dict(connexion.request.get_json())  # noqa: E501
    return 'do some magic!'
