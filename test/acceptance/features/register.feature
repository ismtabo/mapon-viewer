Feature: Http Server

  Background: 
    Given the HTTP endpoint "[CONF:url]"
    Given the mongodb URI "[CONF:mongo.uri]"
    Given the mongodb database "[CONF:mongo.database]"

  Scenario Outline: Server. Register
    Given I generate a UUID and store it in context "register.email"
      And I generate a UUID and store it in context "register.password"
    Given the HTTP path "/auth/register"
      And the HTTP request body with the URL encoded properties
          | email    | [CTXT:register.email]    |
          | password | [CTXT:register.password] |
     When I send a HTTP "POST" request
     Then the HTTP status code must be "200"
      And I search a mongodb single result document in collection "[CONF:mongo.collection]" with filter
          | email    | [CTXT:register.email] |
      And the mongodb search result should be have the following JSON properties
          | email    | [CTXT:register.email] |
      And the HTTP response should contain the cookie "[CONF:session.cookie]"

  Scenario Outline: Server. Register. Missing email
    Given the HTTP path "/auth/register"
      And the HTTP request body with the URL encoded properties
          | password | password |
     When I send a HTTP "POST" request
     Then the HTTP status code must be "400"

  Scenario Outline: Server. Register. Missing email
    Given the HTTP path "/auth/register"
      And the HTTP request body with the URL encoded properties
          | email | email |
     When I send a HTTP "POST" request
     Then the HTTP status code must be "400"

  Scenario Outline: Server. Register. Conflict email
    Given I generate a UUID and store it in context "register.email"
      And I generate a UUID and store it in context "register.password"
      And I generate a salt hash for password "[CTXT:register.password]" and store in context "register.salt"
      And I create mongodb document in collection "[CONF:mongo.collection]" with properties
          | email    | [CTXT:register.email] |
          | password | [CTXT:register.salt]  |
    Given the HTTP path "/auth/register"
      And the HTTP request body with the URL encoded properties
          | email    | [CTXT:register.email]    |
          | password | [CTXT:register.password] |
     When I send a HTTP "POST" request
     Then the HTTP status code must be "409"

  @wip
  Scenario Outline: Server. Register then login
    Given I generate a UUID and store it in context "register.email"
      And I generate a UUID and store it in context "register.password"
    Given the HTTP path "/auth/register"
      And the HTTP request body with the URL encoded properties
          | email    | [CTXT:register.email]    |
          | password | [CTXT:register.password] |
      And I send a HTTP "POST" request
      And the HTTP status code must be "200"
      And the HTTP response should contain the cookie "[CONF:session.cookie]"
      And I keep the HTTP cookies
     When the HTTP path "/auth/login"
      And the HTTP request body with the URL encoded properties
          | email    | [CTXT:register.email]    |
          | password | [CTXT:register.password] |
      And I send a HTTP "POST" request
      And the HTTP status code must be "200"

  Scenario Outline: Server. Register then logout
    Given I generate a UUID and store it in context "register.email"
      And I generate a UUID and store it in context "register.password"
    Given the HTTP path "/auth/register"
      And the HTTP request body with the URL encoded properties
          | email    | [CTXT:register.email]    |
          | password | [CTXT:register.password] |
      And I send a HTTP "POST" request
      And the HTTP status code must be "200"
      And the HTTP response should contain the cookie "[CONF:session.cookie]"
      And I keep the HTTP cookies
     When the HTTP path "/auth/logout"
      And I send a HTTP "POST" request
      And the HTTP status code must be "200"
      And the HTTP response should contain the expired cookie "[CONF:session.cookie]"
