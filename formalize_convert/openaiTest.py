import openai

openai.api_key = ""

response = openai.ChatCompletion.create(
  model="gpt-3.5-turbo",
  messages=[
        {"role": "user", "content": "What is the capital of France?"}
    ]
)

print(response.choices[0].message['content'].strip())