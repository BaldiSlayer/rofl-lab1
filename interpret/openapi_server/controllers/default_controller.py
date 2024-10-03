import connexion
from typing import Dict
from typing import Tuple
from typing import Union

from openapi_server.models.interpret_result import InterpretResult  # noqa: E501
from openapi_server.models.trs import Trs  # noqa: E501
from openapi_server import util


def trs_interpret(trs):  # noqa: E501
    """Check decidability of trs

     # noqa: E501

    :param trs: Trs
    :type trs: dict | bytes

    :rtype: Union[InterpretResult, Tuple[InterpretResult, int], Tuple[InterpretResult, int, Dict[str, str]]
    """
    if connexion.request.is_json:
        trs = Trs.from_dict(connexion.request.get_json())  # noqa: E501
    return 'do some magic!'
