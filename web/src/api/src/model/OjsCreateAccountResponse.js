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
import OjsAccount from './OjsAccount';

/**
 * The OjsCreateAccountResponse model module.
 * @module model/OjsCreateAccountResponse
 * @version version not set
 */
class OjsCreateAccountResponse {
    /**
     * Constructs a new <code>OjsCreateAccountResponse</code>.
     * @alias module:model/OjsCreateAccountResponse
     */
    constructor() { 
        
        OjsCreateAccountResponse.initialize(this);
    }

    /**
     * Initializes the fields of this object.
     * This method is used by the constructors of any subclasses, in order to implement multiple inheritance (mix-ins).
     * Only for internal use.
     */
    static initialize(obj) { 
    }

    /**
     * Constructs a <code>OjsCreateAccountResponse</code> from a plain JavaScript object, optionally creating a new instance.
     * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @param {module:model/OjsCreateAccountResponse} obj Optional instance to populate.
     * @return {module:model/OjsCreateAccountResponse} The populated <code>OjsCreateAccountResponse</code> instance.
     */
    static constructFromObject(data, obj) {
        if (data) {
            obj = obj || new OjsCreateAccountResponse();

            if (data.hasOwnProperty('account')) {
                obj['account'] = OjsAccount.constructFromObject(data['account']);
            }
        }
        return obj;
    }

    /**
     * Validates the JSON data with respect to <code>OjsCreateAccountResponse</code>.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @return {boolean} to indicate whether the JSON data is valid with respect to <code>OjsCreateAccountResponse</code>.
     */
    static validateJSON(data) {
        // validate the optional field `account`
        if (data['account']) { // data not null
          OjsAccount.validateJSON(data['account']);
        }

        return true;
    }


}



/**
 * @member {module:model/OjsAccount} account
 */
OjsCreateAccountResponse.prototype['account'] = undefined;






export default OjsCreateAccountResponse;

