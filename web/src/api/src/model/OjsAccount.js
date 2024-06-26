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
import OjsRole from './OjsRole';

/**
 * The OjsAccount model module.
 * @module model/OjsAccount
 * @version version not set
 */
class OjsAccount {
    /**
     * Constructs a new <code>OjsAccount</code>.
     * @alias module:model/OjsAccount
     */
    constructor() { 
        
        OjsAccount.initialize(this);
    }

    /**
     * Initializes the fields of this object.
     * This method is used by the constructors of any subclasses, in order to implement multiple inheritance (mix-ins).
     * Only for internal use.
     */
    static initialize(obj) { 
    }

    /**
     * Constructs a <code>OjsAccount</code> from a plain JavaScript object, optionally creating a new instance.
     * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @param {module:model/OjsAccount} obj Optional instance to populate.
     * @return {module:model/OjsAccount} The populated <code>OjsAccount</code> instance.
     */
    static constructFromObject(data, obj) {
        if (data) {
            obj = obj || new OjsAccount();

            if (data.hasOwnProperty('id')) {
                obj['id'] = ApiClient.convertToType(data['id'], 'String');
            }
            if (data.hasOwnProperty('name')) {
                obj['name'] = ApiClient.convertToType(data['name'], 'String');
            }
            if (data.hasOwnProperty('role')) {
                obj['role'] = OjsRole.constructFromObject(data['role']);
            }
        }
        return obj;
    }

    /**
     * Validates the JSON data with respect to <code>OjsAccount</code>.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @return {boolean} to indicate whether the JSON data is valid with respect to <code>OjsAccount</code>.
     */
    static validateJSON(data) {
        // ensure the json data is a string
        if (data['id'] && !(typeof data['id'] === 'string' || data['id'] instanceof String)) {
            throw new Error("Expected the field `id` to be a primitive type in the JSON string but got " + data['id']);
        }
        // ensure the json data is a string
        if (data['name'] && !(typeof data['name'] === 'string' || data['name'] instanceof String)) {
            throw new Error("Expected the field `name` to be a primitive type in the JSON string but got " + data['name']);
        }

        return true;
    }


}



/**
 * @member {String} id
 */
OjsAccount.prototype['id'] = undefined;

/**
 * @member {String} name
 */
OjsAccount.prototype['name'] = undefined;

/**
 * @member {module:model/OjsRole} role
 */
OjsAccount.prototype['role'] = undefined;






export default OjsAccount;

