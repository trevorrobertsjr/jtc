---
sidebar_position: 1
---

# Hybrid Cloud?

Hybrid Cloud is the name given to an architecture pattern where a business uses a mix of on-premises resources (i.e. private cloud) and public cloud resources.

They may use data stored on-premises with compute resources in AWS or vice versa.

Here is a high-level diagram to depict this architecture

![Hybrid Cloud](/img/The-Hybrid-Cloud-Model-Figure-1.png)

The bridge may be a virtual private network (VPN), a dedicated network connection like AWS Direct Connect, or even mTLS in a Zero Trust configuration.

In this lab, the Private Cloud will be our local computers where we will run a python script that contains our data pipeline to process data stored in the Public Cloud: An AWS S3 Bucket.