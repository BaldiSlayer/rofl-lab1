#!/usr/bin/env python3

import llm_client
from llm_client.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to http://localhost
# See configuration.py for a list of all supported configuration parameters.
configuration = llm_client.Configuration(
    host="http://llm:8100"
)


# Enter a context with an instance of the API client
with llm_client.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = llm_client.QuestionsApi(api_client)
    get_chat_response_request = llm_client.GetChatResponseRequest.from_dict({
        "prompt": "Привет!",
        "model": "mistral-large-latest",
    })

    try:
        # Api Add Questions
        api_response = api_instance.api_get_chat_response_get_chat_response_post(
            get_chat_response_request)
        print(
            "The response of QuestionsApi->api_get_chat_response_get_chat_response_post:\n")
        pprint(api_response)
    except ApiException as e:
        print("Exception when calling QuestionsApi->api_get_chat_response_get_chat_response_post: %s\n" % e)
