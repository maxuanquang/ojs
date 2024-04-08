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
    instance = new OjsProto.OjsTestCase();
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

  describe('OjsTestCase', function() {
    it('should create an instance of OjsTestCase', function() {
      // uncomment below and update the code to test OjsTestCase
      //var instance = new OjsProto.OjsTestCase();
      //expect(instance).to.be.a(OjsProto.OjsTestCase);
    });

    it('should have the property id (base name: "id")', function() {
      // uncomment below and update the code to test the property id
      //var instance = new OjsProto.OjsTestCase();
      //expect(instance).to.be();
    });

    it('should have the property ofProblemId (base name: "ofProblemId")', function() {
      // uncomment below and update the code to test the property ofProblemId
      //var instance = new OjsProto.OjsTestCase();
      //expect(instance).to.be();
    });

    it('should have the property input (base name: "input")', function() {
      // uncomment below and update the code to test the property input
      //var instance = new OjsProto.OjsTestCase();
      //expect(instance).to.be();
    });

    it('should have the property output (base name: "output")', function() {
      // uncomment below and update the code to test the property output
      //var instance = new OjsProto.OjsTestCase();
      //expect(instance).to.be();
    });

    it('should have the property isHidden (base name: "isHidden")', function() {
      // uncomment below and update the code to test the property isHidden
      //var instance = new OjsProto.OjsTestCase();
      //expect(instance).to.be();
    });

  });

}));
