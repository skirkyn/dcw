# this is just a test lambda
import json


def handler(event, context):
    method = event['requestContext']['http']['method']
    if method != 'POST':
        return {
            'statusCode': 405,
            'body': 'only post is allowed'
        }
    else:
        try:
            body = json.loads(event['body'])
            if body['code'] == 2222:
                return {
                    'statusCode': 200,
                    'body': 'success'
                }
            else:
                return {
                    'statusCode': 500,
                    'body': 'wrong code'
                }
        except Exception as e:
            return {
                'statusCode': 500,
                'body': str(e)
            }
