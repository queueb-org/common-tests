common.queueb.org/tests
=======================

.. contents::
  :local:
  :depth: 2

An extension for Go's standard testing package, ``testing``.

The package provides simple primitives to work with:

- files -- file related helpers (reading, writing, copying, etc).
- golden -- a simple helper for Golden-like tests.
- http -- implements a couple http.Handler routers for testing.
- os -- application arguments override.
- tests -- helpers that might help to test your testing function assets (see files_test.go for example).

It's still in development and used for internal queueb projects.

Installation
------------

Install the package with:

.. code-block:: console

   go get common.queueb.org/tests

Usage
-----

Import the package:

.. code-block:: go

   import "common.queueb.org/tests"

Purpose
-------
Testing package is not a replacement for existent testing-related packages and solutions, it's more or less
'copy-pasting' compatible primitives to remove some chores from developer.

Development
-----------

Documentation
~~~~~~~~~~~~~
Documentation is in development process. Sooner or later it would be available on 
Go Package Documentation: https://pkg.go.dev/common.queueb.org/tests

Source code: https://github.com/queueb-org/common-tests 

Testing
~~~~~~~

.. code-block:: bash

  $ go test -coverprofile=.coverage ./... && go tool cover -func=.coverage

License
-------

See the ``LICENSE`` file in the repository.