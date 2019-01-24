.. Macaroons AuthZ Demo documentation master file, created by
   sphinx-quickstart on Wed Dec 26 10:35:13 2018.
   You can adapt this file completely to your liking, but it should at least
   contain the root `toctree` directive.

Documentation for Macaroon AuthZ Demo project 
=============================================

.. toctree::
   :maxdepth: 2
   :caption: Contents:

   macaroons
   flows



Use Case and Terms
===================

The purpose of this repository is to demonstrate the use of Macaroons as an acceptible authorization scheme for verious projects at the *Centers for Medicare and Medicaid Services* (CMS).
While these principles and concepts in this demo are appliable to a wide variety of use cases, the terminology is focused on specific CMS use cases.
This sections gives an overview of the terminology used in this documentation, as well as the motivating example which inspired this project.

Use Case
--------

We have an API endpoint, which distributes claims history for a given set of Medicare beneficiaries that are assigned to a specific *Accountable Care Organization*.

ACO data can be accessed either by one of their employees, or by the employee of a third-party *vendor* that the ACO has delegated its access to. This means, the API endpoint should only divulge data if one of two conditions are satisfied:

1. The user is an employee of the ACO that it is requesting data for.
2. The user is an employee of a vendor that is authorized by the ACO for which it is retrieving data.

Terms
-----




Indices and tables
==================

* :ref:`genindex`
* :ref:`modindex`
* :ref:`search`
