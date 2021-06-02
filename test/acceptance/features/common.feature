Feature: Http Server

  Background: 
    Given the HTTP endpoint "[CONF:url]"

  Scenario Outline: Server. Login - Unauthorized method
    Given the HTTP path "/auth/login"
     When I send a HTTP "<method>" request
     Then the HTTP status code must be "405"
        
    Examples: method: <method>
          | method      |
          | GET         |
          | PUT         |
          | PATCH       |
          | DELETE      |

  Scenario Outline: Server. Logout - Unauthorized method
    Given the HTTP path "/auth/logout"
     When I send a HTTP "<method>" request
     Then the HTTP status code must be "405"
        
    Examples: method: <method>
          | method      |
          | GET         |
          | PUT         |
          | PATCH       |
          | DELETE      |

  Scenario Outline: Server. Register - Unauthorized method
    Given the HTTP path "/auth/register"
     When I send a HTTP "<method>" request
     Then the HTTP status code must be "405"
        
    Examples: method: <method>
          | method      |
          | GET         |
          | PUT         |
          | PATCH       |
          | DELETE      |