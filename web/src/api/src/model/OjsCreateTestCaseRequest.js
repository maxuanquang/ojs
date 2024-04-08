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

import ApiClient from '../ApiClient';

/**
 * The OjsCreateTestCaseRequest model module.
 * @module model/OjsCreateTestCaseRequest
 * @version version not set
 */
class OjsCreateTestCaseRequest {
    /**
     * Constructs a new <code>OjsCreateTestCaseRequest</code>.
     * @alias module:model/OjsCreateTestCaseRequest
     */
    constructor() { 
        
        OjsCreateTestCaseRequest.initialize(this);
    }

    /**
     * Initializes the fields of this object.
     * This method is used by the constructors of any subclasses, in order to implement multiple inheritance (mix-ins).
     * Only for internal use.
     */
    static initialize(obj) { 
    }

    /**
     * Constructs a <code>OjsCreateTestCaseRequest</code> from a plain JavaScript object, optionally creating a new instance.
     * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @param {module:model/OjsCreateTestCaseRequest} obj Optional instance to populate.
     * @return {module:model/OjsCreateTestCaseRequest} The populated <code>OjsCreateTestCaseRequest</code> instance.
     */
    static constructFromObject(data, obj) {
        if (data) {
            obj = obj || new OjsCreateTestCaseRequest();

            if (data.hasOwnProperty('ofProblemId')) {
                obj['ofProblemId'] = ApiClient.convertToType(data['ofProblemId'], 'String');
            }
            if (data.hasOwnProperty('input')) {
                obj['input'] = ApiClient.convertToType(data['input'], 'String');
            }
            if (data.hasOwnProperty('output')) {
                obj['output'] = ApiClient.convertToType(data['output'], 'String');
            }
            if (data.hasOwnProperty('isHidden')) {
                obj['isHidden'] = ApiClient.convertToType(data['isHidden'], 'Boolean');
            }
        }
        return obj;
    }

    /**
     * Validates the JSON data with respect to <code>OjsCreateTestCaseRequest</code>.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @return {boolean} to indicate whether the JSON data is valid with respect to <code>OjsCreateTestCaseRequest</code>.
     */
    static validateJSON(data) {
        // ensure the json data is a string
        if (data['ofProblemId'] && !(typeof data['ofProblemId'] === 'string' || data['ofProblemId'] instanceof String)) {
            throw new Error("Expected the field `ofProblemId` to be a primitive type in the JSON string but got " + data['ofProblemId']);
        }
        // ensure the json data is a string
        if (data['input'] && !(typeof data['input'] === 'string' || data['input'] instanceof String)) {
            throw new Error("Expected the field `input` to be a primitive type in the JSON string but got " + data['input']);
        }
        // ensure the json data is a string
        if (data['output'] && !(typeof data['output'] === 'string' || data['output'] instanceof String)) {
            throw new Error("Expected the field `output` to be a primitive type in the JSON string but got " + data['output']);
        }

        return true;
    }


}



/**
 * @member {String} ofProblemId
 */
OjsCreateTestCaseRequest.prototype['ofProblemId'] = undefined;

/**
 * @member {String} input
 */
OjsCreateTestCaseRequest.prototype['input'] = undefined;

/**
 * @member {String} output
 */
OjsCreateTestCaseRequest.prototype['output'] = undefined;

/**
 * @member {Boolean} isHidden
 */
OjsCreateTestCaseRequest.prototype['isHidden'] = undefined;






export default OjsCreateTestCaseRequest;

