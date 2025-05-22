# Problem Statement

During their trip to Bangkok, a group of friends collectively incurred various expenses, including flight tickets, taxi fares, museum entry fees, and meals. Each expense was covered spontaneously by different individuals within the group.

They now require an application that can calculate each individual's balance of owed money and payments within the group, detailing how much each person needs to reimburse and how much others owe them. Also, there is a need to simplify debts within groups so that they have the **least number of transactions** taking place between different group members. For Example 

A -> paid -> B   -  $40
B -> paid -> C   -  $40
C -> paid -> A   -  $10

Rather than each individual paying their debt back, we can simplify things and have only one transaction from C to A ($30), which will settle everyone's balance.

You need to create a Web Application with the necessary Rest Endpoints. 

# What do we expect?

- An application backend that fulfills the requirements mentioned above using rest API. No UI is needed.
- The submission must be pushed to this repository under the `main` branch.
- Instructions on how to run the application.
- If we move forward to the following interview, we'll ask for a short (up to 5 minutes) live demo of the application showcasing its features.

# What do we like in a Solution?

- Good software design, following necessary design practices.
- Proper documentation for APIs and the overall application. The reviewers should clearly understand the code structure, how to run the service and how to use the API.
- Functionality should be modularized and distributed across multiple classes, methods, or files to avoid creating a single overly complex class, method, or file, commonly referred to as a "God class".
- Code should be easily extensible & maintainable.

# Architecture Review

The second exercise is creating a sample architecture for your proposed solution. Let's suppose you're developing a SaaS product with the capabilities above.

For this exercise, prepare one slide explaining the solution design in the cloud of your choice. We value topics such as scalability, reliability, portability, etc.












