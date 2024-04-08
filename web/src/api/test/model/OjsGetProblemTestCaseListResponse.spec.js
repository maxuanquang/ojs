/**
 * ojs.proto
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: version not set
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 *
 */

(function(root, factory) {
  if (typeof define === 'function' && define.amd) {
    // AMD.
    define(['expect.js', process.cwd()+'/src/index'], factory);
  } else if (typeof module === 'object' && module.exports) {
    // CommonJS-like environments that support module.exports, like Node.
    factory(require('expect.js'), require(process.cwd()+'/src/index'));
  } else {
    // Browser globals (root is window)
    factory(root.expect, root.OjsProto);
  }
}(this, function(expect, OjsProto) {
  'use strict';

  var instance;

  beforeEach(function() {
    instance = new OjsProto.OjsGetProblemTestCaseListResponse();
  });

  var getProperty = function(object, getter, property) {
    // Use getter method if present; otherwise, get the property directly.
    if (typeof object[getter] === 'function')
      return object[getter]();
    else
      return object[property];
  }

  var setProperty = function(object, setter, property, value) {
    // Use setter method if present; otherwise, set the property directly.
    if (typeof object[setter] === 'function')
      object[setter](value);
    else
      object[property] = value;
  }

  describe('OjsGetProblemTestCaseListResponse', function() {
    it('should create an instance of OjsGetProblemTestCaseListResponse', function() {
      // uncomment below and update the code to test OjsGetProblemTestCaseListResponse
      //var instance = new OjsProto.OjsGetProblemTestCaseListResponse();
      //expect(instance).to.be.a(OjsProto.OjsGetProblemTestCaseListResponse);
    });

    it('should have the property testCases (base name: "testCases")', function() {
      // uncomment below and update the code to test the property testCases
      //var instance = new OjsProto.OjsGetProblemTestCaseListResponse();
      //expect(instance).to.be();
    });

    it('should have the property totalTestCasesCount (base name: "totalTestCasesCount")', function() {
      // uncomment below and update the code to test the property totalTestCasesCount
      //var instance = new OjsProto.OjsGetProblemTestCaseListResponse();
      //expect(instance).to.be();
    });

  });

}));
