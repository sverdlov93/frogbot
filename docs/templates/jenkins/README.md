[Go back to the Installation documentation page](../../../README.md)

# Set up Frogbot Using Jenkins

### 🖥️ Follow these steps to install Frogbot on Jenkins

<details>
  <summary>1️⃣ Install Jenkins 'Generic Webhook Trigger' plugin </summary>

From your Jenkins dashboard navigate to **Manage Jenkins** > **Manage Plugins** and select the **Available** tab.
Use the search bar to find **Generic Webhook Trigger** ([more info](https://plugins.jenkins.io/generic-webhook-trigger/)).

</details>

---
<details>
  <summary>2️⃣ Connect the Webhook on your Git provider </summary>

<details>
      <summary>Bitbucket Server</summary>

- Webhook URL: `JENKINS_URL/generic-webhook-trigger/invoke`
- Go to repository settings, select Webhooks, and create a new webhook.
  <img src="../../../images/bitbucket-webhook-setup.png">
- Set the webhook URL `https://jenkinsUrl/generic-webhook-trigger/invoke`
  <img src="../../../images/bitbucketserver-create-webhook.png">
</details>

<details>
    <summary>GitHub</summary>

- Webhook URL: `JENKINS_URL/generic-webhook-trigger/invoke`
- Go to repository settings and create a new webhook:
  <img src="../../../images/github-new-webhook.png">

- Add a new webhook:
  <img src="../../../images/github-webhook-setup.png">

- Set up trigger:
  <img src="../../../images/github-trigger-event.png">

</details>

<details>
  <summary>Azure Repos</summary>

- Webhook URL: `JENKINS_URL/generic-webhook-trigger/invoke`
- [Set Up Azure Repos Jenkins Webhook](https://learn.microsoft.com/en-us/azure/devops/service-hooks/services/jenkins?view=azure-devops)

</details>

<details>
   <summary>GitLab</summary>

- Go to your project settings and select webhooks.
- Set up a webhook with merge request events.
- Fill in the URL: `JENKINS URL/generic-webhook-trigger/invoke`
  <img src="../../../images/GitLab_webhook.png">

</details>

</details>

---
<details>
  <summary>3️⃣ Optional - setting JobToken</summary>

  - When using the plugin in several jobs, you will have the same URL trigger all jobs. If you
    want to trigger only a certain job you can use the **JobToken** in the URL to specify what job needs to be executed.
  - Webhook URL with **JobToken** : `JENKINS_URL/generic-webhook-trigger/invoke?token=MyJobToken`
  - On some Git providers the JobToken is called Secret Token.
  - Read more [JobToken Docs](https://plugins.jenkins.io/generic-webhook-trigger/#plugin-content-trigger-only-specific-job)
</details>

---
<details>
  <summary>4️⃣ Set up credentials</summary>

- Set up the following credentials using Jenkins credentials functionality, as **Secret Text**:
    - **JF_URL** - JFrog Platform URL (Example: "https://acme.jfrog.io")
    - **JF_ACCESS_TOKEN** *or* **JF_USER** & **JF_PASSWORD** - JFrog Credentials
    - **JF_GIT_TOKEN** - access token with read&write access to the Git repository
- [How to use credentials with Jenkins](https://www.jenkins.io/doc/book/using/using-credentials/)

</details>

---
<details>
  <summary>5️⃣ Prepare Jenkins Agent</summary>

- It is essential to have the appropriate package manager used by the scanned project installed on the Jenkins Agent. For instance, if the project uses an npm project, you need to have the npm client installed.

</details>

---
<details>
  <summary>6️⃣ Scanning pull requests</summary>

Create a new pipeline job using [this](./scan-pull-request.jenkinsfile) Jenkinsfile template.
  
<img src="../../../images/jenkins-pipeline-select.png" width="650"> 

Enable the ‘Generic Webhook Trigger’:

<img src="../../../images/jenkins-build-trigger.png">

And paste it here : 

<img src="../../../images/jenkins-paste-pipeline.png" width="650"> 

</details>

---
<details>
  <summary>7️⃣ Scanning repository branches and fixing issues</summary>

  
Create a new pipeline job using [this](./scan-repository.jenkinsfile) Jenkinsfile template.

<img src="../../../images/jenkins-pipeline-select.png" width="650">  

And paste it here :

<img src="../../../images/jenkins-paste-pipeline.png" width="650"> 


</details>
