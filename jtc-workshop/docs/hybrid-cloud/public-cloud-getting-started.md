---
sidebar_position: 2
---

# Setting Up Our Data in the Public Cloud with Amazon S3

:::info

I designed this workshop to run properly in the AWS us-east-2 and GCP us-central1 regions. You are free to choose a different region, but you may need to do some alterations to my documented steps to complete the labs.

:::

First, let's get our data file that we will use for the workshop: <a href="/files/nyc_reviews.csv" download>Download NYC Reviews CSV</a>

The data file contains 1000 reviews of New York City that I generated using ChatGPT. If you are interested in working with data on a larger scale, you can always check sites like [Kaggle](https://www.kaggle.com/) and listings of [free public APIs](https://github.com/public-apis/public-apis).

Next, let us create an Amazon S3 bucket and put our data file there.

Login to the AWS Console. You should see something similar to the following:
![AWS Console](/img/s3-00.png)

In the `Search` bar in the top left, type `s3`. The S3 service should come up in a dropdown.
![AWS Console](/img/s3-01.png)

On the Amazon S3 welcome screen, click `Create Bucket`
![AWS Console](/img/s3-02.png)

Provide a very unique bucket name. S3 is a global service, and each bucket name is registered in DNS. So, it is important you select a name that no one else has...No pressure! 😅
![AWS Console](/img/s3-03.png)

Click on your bucket name so that we can upload our file.
![AWS Console](/img/s3-04.png)

Click `Upload`
![AWS Console](/img/s3-05.png)

You can either drag the nyc_reviews.csv file to the window or click `Add files`
![AWS Console](/img/s3-06.png)

Now that our data is uploaded to Amazon S3, let's run our pipeline to generate our **clean** data based on our Data Scientists' requests.