![build](https://github.com/mathisve/latest-youtube-videos-lambda/actions/workflows/go.yaml/badge.svg)

# latest-youtube-videos-lambda
Serverless Lambda function to get the latest Youtube videos of a channel.

The Youtube data API requires you to use an API key. If you have a serverless website with no backend, you can't use this API unless you wrap the API in a Lambda function behind an API Gateway.
This is what I did here.

## Env vars
You need to set these environment variables in the AWS Lambda console! **Ohterwise it will not work!**
Get an API key [here](https://console.cloud.google.com/apis/credentials)! Remember that you need to Activate the Youtube data API too!

![environment variables](https://raw.githubusercontent.com/mathisve/latest-youtube-videos-lambda/master/img/youtubeapi.png)

## Caching
This function caches the response from the Youtube API for **15 minutes**.
If more than 15 minutes have passed since the last request, it will query the youtube API for fresh data.
This can be changed by altering the `cacheTime` constant.

## Results
By default it will only query for the **10 latest videos**.
This too can be altered by changing the `maxResults` constant.

## Thank you!!
