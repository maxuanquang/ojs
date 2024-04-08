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
import OjsTestCase from './OjsTestCase';

/**
 * The OjsGetProblemTestCaseListResponse model module.
 * @module model/OjsGetProblemTestCaseListResponse
 * @version version not set
 */
class OjsGetProblemTestCaseListResponse {
    /**
     * Constructs a new <code>OjsGetProblemTestCaseListResponse</code>.
     * @alias module:model/OjsGetProblemTestCaseListResponse
     */
    constructor() { 
        
        OjsGetProblemTestCaseListResponse.initialize(this);
    }

    /**
     * Initializes the fields of this object.
     * This method is used by the constructors of any subclasses, in order to implement multiple inheritance (mix-ins).
     * Only for internal use.
     */
    static initialize(obj) { 
    }

    /**
     * Constructs a <code>OjsGetProblemTestCaseListResponse</code> from a plain JavaScript object, optionally creating a new instance.
     * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @param {module:model/OjsGetProblemTestCaseListResponse} obj Optional instance to populate.
     * @return {module:model/OjsGetProblemTestCaseListResponse} The populated <code>OjsGetProblemTestCaseListResponse</code> instance.
     */
    static constructFromObject(data, obj) {
        if (data) {
            obj = obj || new OjsGetProblemTestCaseListResponse();

            if (data.hasOwnProperty('testCases')) {
                obj['testCases'] = ApiClient.convertToType(data['testCases'], [OjsTestCase]);
            }
            if (data.hasOwnProperty('totalTestCasesCount')) {
                obj['totalTestCasesCount'] = ApiClient.convertToType(data['totalTestCasesCount'], 'String');
            }
        }
        return obj;
    }

    /**
     * Validates the JSON data with respect to <code>OjsGetProblemTestCaseListResponse</code>.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @return {boolean} to indicate whether the JSON data is valid with respect to <code>OjsGetProblemTestCaseListResponse</code>.
     */
    static validateJSON(data) {
        if (data['testCases']) { // data not null
            // ensure the json data is an array
            if (!Array.isArray(data['testCases'])) {
                throw new Error("Expected the field `testCases` to be an array in the JSON data but got " + data['testCases']);
            }
            // validate the optional field `testCases` (array)
            for (const item of data['testCases']) {
                OjsTestCase.validateJSON(item);
            };
        }
        // ensure the json data is a string
        if (data['totalTestCasesCount'] && !(typeof data['totalTestCasesCount'] === 'string' || data['totalTestCasesCount'] instanceof String)) {
            throw new Error("Expected the field `totalTestCasesCount` to be a primitive type in the JSON string but got " + data['totalTestCasesCount']);
        }

        return true;
    }


}



/**
 * @member {Array.<module:model/OjsTestCase>} testCases
 */
OjsGetProblemTestCaseListResponse.prototype['testCases'] = undefined;

/**
 * @member {String} totalTestCasesCount
 */
OjsGetProblemTestCaseListResponse.prototype['totalTestCasesCount'] = undefined;






export default OjsGetProblemTestCaseListResponse;
