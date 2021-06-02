Feature: Http Server

  Background: 
    Given the HTTP endpoint "[CONF:url]"
    Given the mongodb URI "[CONF:mongo.uri]"
    Given the mongodb database "[CONF:mongo.database]"

  Scenario Outline: Server. Login
    Given I generate a UUID and store it in context "login.email"
      And I generate a UUID and store it in context "login.password"
      And I generate a salt hash for password "[CTXT:login.password]" and store in context "login.salt"
      And I create mongodb document in collection "[CONF:mongo.collection]" with properties
          | email    | [CTXT:login.email] |
          | password | [CTXT:login.salt]  |
    Given the HTTP path "/auth/login"
      And the HTTP request body with the URL encoded properties
          | email    | [CTXT:login.email]    |
          | password | [CTXT:login.password] |
     When I send a HTTP "POST" request
     Then the HTTP status code must be "200"
      And the HTTP response should contain the cookie "[CONF:session.cookie]"

  Scenario Outline: Server. Login. Missing email
    Given the HTTP path "/auth/login"
      And the HTTP request body with the URL encoded properties
          | password | password |
     When I send a HTTP "POST" request
     Then the HTTP status code must be "400"

  Scenario Outline: Server. Login. Missing email
    Given the HTTP path "/auth/login"
      And the HTTP request body with the URL encoded properties
          | email | email |
     When I send a HTTP "POST" request
     Then the HTTP status code must be "400"

  Scenario Outline: Server. Login. Unknown email
    Given I generate a UUID and store it in context "login.email"
      And I generate a UUID and store it in context "login.password"
      And I generate a salt hash for password "[CTXT:login.password]" and store in context "login.salt"
      And I create mongodb document in collection "[CONF:mongo.collection]" with properties
          | email    | [CTXT:login.email] |
          | password | [CTXT:login.salt]  |
    Given the HTTP path "/auth/login"
      And the HTTP request body with the URL encoded properties
          | email    | unknown               |
          | password | [CTXT:login.password] |
     When I send a HTTP "POST" request
     Then the HTTP status code must be "401"

  Scenario Outline: Server. Login. Invalid password
    Given I generate a UUID and store it in context "login.email"
      And I generate a UUID and store it in context "login.password"
      And I generate a salt hash for password "[CTXT:login.password]" and store in context "login.salt"
      And I create mongodb document in collection "[CONF:mongo.collection]" with properties
          | email    | [CTXT:login.email] |
          | password | [CTXT:login.salt]  |
    Given the HTTP path "/auth/login"
      And the HTTP request body with the URL encoded properties
          | email    | [CTXT:login.email] |
          | password | unknown            |
     When I send a HTTP "POST" request
     Then the HTTP status code must be "401"

  Scenario Outline: Server. Login then logout
    Given I generate a UUID and store it in context "login.email"
      And I generate a UUID and store it in context "login.password"
      And I generate a salt hash for password "[CTXT:login.password]" and store in context "login.salt"
      And I create mongodb document in collection "[CONF:mongo.collection]" with properties
          | email    | [CTXT:login.email] |
          | password | [CTXT:login.salt]  |
    Given the HTTP path "/auth/login"
      And the HTTP request body with the URL encoded properties
          | email    | [CTXT:login.email]    |
          | password | [CTXT:login.password] |
      And I send a HTTP "POST" request
      And the HTTP status code must be "200"
      And I keep the HTTP cookies
     When the HTTP path "/auth/logout"
      And I send a HTTP "POST" request
      And the HTTP response should contain the expired cookie "[CONF:session.cookie]"
