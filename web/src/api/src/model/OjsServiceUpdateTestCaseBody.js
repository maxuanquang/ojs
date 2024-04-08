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
 * The OjsServiceUpdateTestCaseBody model module.
 * @module model/OjsServiceUpdateTestCaseBody
 * @version version not set
 */
class OjsServiceUpdateTestCaseBody {
    /**
     * Constructs a new <code>OjsServiceUpdateTestCaseBody</code>.
     * @alias module:model/OjsServiceUpdateTestCaseBody
     */
    constructor() { 
        
        OjsServiceUpdateTestCaseBody.initialize(this);
    }

    /**
     * Initializes the fields of this object.
     * This method is used by the constructors of any subclasses, in order to implement multiple inheritance (mix-ins).
     * Only for internal use.
     */
    static initialize(obj) { 
    }

    /**
     * Constructs a <code>OjsServiceUpdateTestCaseBody</code> from a plain JavaScript object, optionally creating a new instance.
     * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @param {module:model/OjsServiceUpdateTestCaseBody} obj Optional instance to populate.
     * @return {module:model/OjsServiceUpdateTestCaseBody} The populated <code>OjsServiceUpdateTestCaseBody</code> instance.
     */
    static constructFromObject(data, obj) {
        if (data) {
            obj = obj || new OjsServiceUpdateTestCaseBody();

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
     * Validates the JSON data with respect to <code>OjsServiceUpdateTestCaseBody</code>.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @return {boolean} to indicate whether the JSON data is valid with respect to <code>OjsServiceUpdateTestCaseBody</code>.
     */
    static validateJSON(data) {
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
 * @member {String} input
 */
OjsServiceUpdateTestCaseBody.prototype['input'] = undefined;

/**
 * @member {String} output
 */
OjsServiceUpdateTestCaseBody.prototype['output'] = undefined;

/**
 * @member {Boolean} isHidden
 */
OjsServiceUpdateTestCaseBody.prototype['isHidden'] = undefined;






export default OjsServiceUpdateTestCaseBody;

