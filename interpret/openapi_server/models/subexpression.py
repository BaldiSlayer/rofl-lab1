from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from openapi_server.models.base_model import Model
from openapi_server.models.letter import Letter
from openapi_server import util

from openapi_server.models.letter import Letter  # noqa: E501


class Subexpression(Model):
    """NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).

    Do not edit the class manually.
    """

    def __init__(self, letter=None, args=None):  # noqa: E501
        """Subexpression - a model defined in OpenAPI

        :param letter: The letter of this Subexpression.  # noqa: E501
        :type letter: Letter
        :param args: The args of this Subexpression.  # noqa: E501
        :type args: List[object]
        """
        self.openapi_types = {
            'letter': Letter,
            'args': List[object]
        }

        self.attribute_map = {
            'letter': 'letter',
            'args': 'args'
        }

        self._letter = letter
        self._args = args

    @classmethod
    def from_dict(cls, dikt) -> 'Subexpression':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The Subexpression of this Subexpression.  # noqa: E501
        :rtype: Subexpression
        """
        return util.deserialize_model(dikt, cls)

    @property
    def letter(self) -> Letter:
        """Gets the letter of this Subexpression.


        :return: The letter of this Subexpression.
        :rtype: Letter
        """
        return self._letter

    @letter.setter
    def letter(self, letter: Letter):
        """Sets the letter of this Subexpression.


        :param letter: The letter of this Subexpression.
        :type letter: Letter
        """
        if letter is None:
            raise ValueError("Invalid value for `letter`, must not be `None`")  # noqa: E501

        self._letter = letter

    @property
    def args(self) -> List[object]:
        """Gets the args of this Subexpression.


        :return: The args of this Subexpression.
        :rtype: List[object]
        """
        return self._args

    @args.setter
    def args(self, args: List[object]):
        """Sets the args of this Subexpression.


        :param args: The args of this Subexpression.
        :type args: List[object]
        """

        self._args = args