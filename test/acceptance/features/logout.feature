Feature: Http Server

  Background: 
    Given the HTTP endpoint "[CONF:url]"
    Given the mongodb URI "[CONF:mongo.uri]"
    Given the mongodb database "[CONF:mongo.database]"

  Scenario Outline: Server. Logout
    Given the HTTP path "/auth/logout"
     When I send a HTTP "POST" request
     Then the HTTP status code must be "200"
