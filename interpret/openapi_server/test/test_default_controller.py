import unittest

from flask import json

from openapi_server.models.interpret_result import InterpretResult  # noqa: E501
from openapi_server.models.trs import Trs  # noqa: E501
from openapi_server.test import BaseTestCase


class TestDefaultController(BaseTestCase):
    """DefaultController integration test stubs"""

    def test_trs_interpret(self):
        """Test case for trs_interpret

        Check decidability of trs
        """
        trs = {"variables": [null, null], "rules": [{"lhs": {"args": [null, null], "letter": {"name": "name", "isVariable": True}}, "rhs": {"args": [null, null], "letter": {"name": "name", "isVariable": True}}}, {"lhs": {"args": [null, null], "letter": {"name": "name", "isVariable": True}}, "rhs": {"args": [null, null], "letter": {"name": "name", "isVariable": True}}}], "interpretations": [
            {"args": [null, null], "monomials": [{"variable": "variable", "coefficient": 1, "power": 1}, {"variable": "variable", "coefficient": 1, "power": 1}], "name": "name", "constants": [1, 1]}, {"args": [null, null], "monomials": [{"variable": "variable", "coefficient": 1, "power": 1}, {"variable": "variable", "coefficient": 1, "power": 1}], "name": "name", "constants": [1, 1]}]}
        headers = {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
        }
        response = self.client.open(
            '/trs/interpret',
            method='POST',
            headers=headers,
            data=json.dumps(trs),
            content_type='application/json')
        self.assert200(response,
                       'Response body is : ' + response.data.decode('utf-8'))


if __name__ == '__main__':
    unittest.main()
