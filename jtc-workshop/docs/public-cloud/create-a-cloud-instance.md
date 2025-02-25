---
sidebar_position: 1
---

# Create an Amazon EC2 Instance

Give it a **Name** ex: my_instance.
Select the **OS** Ubuntu 24.04 LTS
Select the **Architecture** x86
Select the **Instance type** t3.small - this is not the free tier instance type, but we need 2 vCPUs for our lab, and this instance type only has 1 vCPU
For the **Key pair** click `Create new key pair` unless you have an existing one in the dropdown box (select .pem unless you're using PuTTY on Windows). If you created a key pair, the private key pem/ppk file will automatically download.
In **Network settings**, make sure `Allow SSH traffic from` is **checked** and in the dropdown box, select **My IP**
Finally, click **Launch instance**

Now, as developers, who enjoyed all that pointing and clicking?...perhaps, there is a better way!

:::danger

ALWAYS delete your test EC2 instances, especially if they have public IP addresses.

:::

Add **Markdown or React** files to `src/pages` to create a **standalone page**:

- `src/pages/index.js` → `localhost:3000/`
- `src/pages/foo.md` → `localhost:3000/foo`
- `src/pages/foo/bar.js` → `localhost:3000/foo/bar`

## Create your first React Page

Create a file at `src/pages/my-react-page.js`:

```jsx title="src/pages/my-react-page.js"
import React from 'react';
import Layout from '@theme/Layout';

export default function MyReactPage() {
  return (
    <Layout>
      <h1>My React page</h1>
      <p>This is a React page</p>
    </Layout>
  );
}
```

A new page is now available at [http://localhost:3000/my-react-page](http://localhost:3000/my-react-page).

## Create your first Markdown Page

Create a file at `src/pages/my-markdown-page.md`:

```mdx title="src/pages/my-markdown-page.md"
# My Markdown page

This is a Markdown page
```

A new page is now available at [http://localhost:3000/my-markdown-page](http://localhost:3000/my-markdown-page).
