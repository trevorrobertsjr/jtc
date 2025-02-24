---
sidebar_position: 2
---

# Setting up our Public Cloud on Amazon S3

:::info

I designed this workshop to run properly in the AWS us-east-1 and GCP us-central1 regions. You are free to choose a different region, but you may need to do some alterations to my documented steps to complete the labs.

:::

Let us create an Amazon S3 bucket and put our data file there.

First, let's get our data file that we will use for the workshop: <a href="/files/nyc_reviews.csv" download>Download NYC Reviews CSV</a>

The data file contains 1000 reviews of New York City that I generated using ChatGPT. If you are interested in working with real data on a larger scale, you can always check sites like [Kaggle](https://www.kaggle.com/) and listings of [free public APIs](https://github.com/public-apis/public-apis).

Now that we have our data, let's host it on AWS

