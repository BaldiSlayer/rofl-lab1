# ProcessQuestionsResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**result** | **List[str]** |  |

## Example

```python
from openapi_client.models.process_questions_response import ProcessQuestionsResponse

# TODO update the JSON string below
json = "{}"
# create an instance of ProcessQuestionsResponse from a JSON string
process_questions_response_instance = ProcessQuestionsResponse.from_json(json)
# print the JSON string representation of the object
print(ProcessQuestionsResponse.to_json())

# convert the object into a dict
process_questions_response_dict = process_questions_response_instance.to_dict()
# create an instance of ProcessQuestionsResponse from a dict
process_questions_response_from_dict = ProcessQuestionsResponse.from_dict(process_questions_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
