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
* Enum class OjsSubmissionStatus.
* @enum {}
* @readonly
*/
export default class OjsSubmissionStatus {
    
        /**
         * value: "UndefinedStatus"
         * @const
         */
        "UndefinedStatus" = "UndefinedStatus";

    
        /**
         * value: "Submitted"
         * @const
         */
        "Submitted" = "Submitted";

    
        /**
         * value: "Executing"
         * @const
         */
        "Executing" = "Executing";

    
        /**
         * value: "Finished"
         * @const
         */
        "Finished" = "Finished";

    

    /**
    * Returns a <code>OjsSubmissionStatus</code> enum value from a Javascript object name.
    * @param {Object} data The plain JavaScript object containing the name of the enum value.
    * @return {module:model/OjsSubmissionStatus} The enum <code>OjsSubmissionStatus</code> value.
    */
    static constructFromObject(object) {
        return object;
    }
}

