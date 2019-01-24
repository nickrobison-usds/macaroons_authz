=========================
User Authentication Flows
=========================
:Date: 2019-01-18
:Version: 0.0.1
:Authors: - Nick Robison

This document outlines two example user authentication flows that illustrate various ways that macaroons can be used to handle user authentication and authorization.

The flows can be executed via the *cli* application, which allows the user to modify various parts of the authentication flow and view the various responses.
The build instructions are given in the main `README <../../README.md#cli-client>`_.

Dynamic Flow
============

The *dynamic* flow represents an authentication method in which the service returning the data generates and manages its own macaroons.
These macaroons require fulfillment of a series of *third-party caveats* due to the fact that the API service has no information as to the relation between entities (ACOS and vendors) and various users.

Running the flow
----------------

This flow assumes that the application is running and the user has logged into the system.
See the login section for more details.
In addition, the CLI client should be built and available in the user path.

For this flow, we want to verify that a user can successfully access data for a given ACO provided that either they are employed by the ACO, or are employed by a vendor that is associated to that ACO.

*Step 1: Make initial request*

Run the CLI client with the following command:

.. code-block:: bash

   ./cli "Test User 1" "Test ACO 1"

The command should fail, because *Test User 1* is not assigned to *Test ACO 1*.
In order to resolve the authorization error, we need to assign the user to ACO, and retry the request.

.. image:: images/flows/user-failure.png
    :width: 500
    :alt: Unsuccesful user request


*Step 2: Assign user to ACO*

Assigining a user can be done via the web application.
By default, the application runs at the address `http://localhost:8080`.
You'll need to initially login, details on how to do so are given in the login section.
From there, navigate to the `Users` tab.

.. image:: images/flows/user-screen.png
    :width: 500
    :alt: Application page showing all registered users

Select the *assign* option in the user table, and select `ACO` for the entity type and `Test ACO 1` for the `Entity`, then click *Assign*.

.. image:: images/flows/user-assign.png
    :width: 500
    :alt: Assign user to ACO

Now, you should be able to execute the CLI command again and see that the request succeeds!

.. code-block:: bash

   ./cli "Test User 1" "Test ACO 1"

.. image:: images/flows/user-success.png
    :width: 500
    :alt: Successful user request

*Step 3: Request data as a vendor*

Run as a vendor user

Fail

*Step 4: Assign user to vendor*

Add the user to the vendor

Fails

*Step 5: Assign vendor to ACO*

Add the vendor to the ACO, success!


Delegated Flow
==============





Logging In
==========



|Login|

