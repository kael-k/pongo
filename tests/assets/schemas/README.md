## How to add new tests
Every folder here is tested in various part of the test files, but the main one that validate the content is `test/assets_test.go`: this test validates the following checks on assets:
* a `pongoschema.json` file is present and a valid json-marshalled PonGO schema
* a `jsonschema.json` file is present and is a valid JSON schema
* **at least one** file with name `^ok.*\.json$` which contains a valid JSON which is correctly validated by both the PonGO schema and JSON schema
* **at least one** file with name `^ko.*\.json$` which contains a valid JSON which is not validated by both the PonGO schema and JSON schema
* optionally a `README.md` which give details about the test
* if there is any file that is not a part of the before said points, the test will fail