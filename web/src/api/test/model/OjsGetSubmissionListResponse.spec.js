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
    instance = new OjsProto.OjsGetSubmissionListResponse();
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

  describe('OjsGetSubmissionListResponse', function() {
    it('should create an instance of OjsGetSubmissionListResponse', function() {
      // uncomment below and update the code to test OjsGetSubmissionListResponse
      //var instance = new OjsProto.OjsGetSubmissionListResponse();
      //expect(instance).to.be.a(OjsProto.OjsGetSubmissionListResponse);
    });

    it('should have the property submissions (base name: "submissions")', function() {
      // uncomment below and update the code to test the property submissions
      //var instance = new OjsProto.OjsGetSubmissionListResponse();
      //expect(instance).to.be();
    });

    it('should have the property totalSubmissionsCount (base name: "totalSubmissionsCount")', function() {
      // uncomment below and update the code to test the property totalSubmissionsCount
      //var instance = new OjsProto.OjsGetSubmissionListResponse();
      //expect(instance).to.be();
    });

  });

}));
