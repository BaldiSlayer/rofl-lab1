#!/usr/bin/env python3

import llm_client.openapi_client as openapi_client
from llm_client.openapi_client.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to http://localhost
# See configuration.py for a list of all supported configuration parameters.
configuration = openapi_client.Configuration(
    host="http://llm:8100"
)


# Enter a context with an instance of the API client
with openapi_client.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = openapi_client.QuestionsApi(api_client)
    get_chat_response_request = openapi_client.GetChatResponseRequest()
    get_chat_response_request.prompt = "Привет!"

    try:
        # Api Add Questions
        api_response = api_instance.api_get_chat_response_get_chat_response_post(
            get_chat_response_request)
        print(
            "The response of QuestionsApi->api_get_chat_response_get_chat_response_post:\n")
        pprint(api_response)
    except ApiException as e:
        print("Exception when calling QuestionsApi->api_get_chat_response_get_chat_response_post: %s\n" % e)
