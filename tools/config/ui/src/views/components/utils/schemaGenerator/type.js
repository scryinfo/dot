'use strict';

let isObject = function (value) {
    return (null !== value && typeof value === typeof {} && !isArray(value));
};

let isNumber = function (value) {
    return !isArray(value) && (value - parseFloat(value) + 1) >= 0;
};

let isArray = function (value) {
    return (value instanceof Array);
};

let isString = function (value) {
    return (typeof value === typeof '');
};

let isNull = function (value) {
    return (null === value);
};

let isBoolean = function (value) {
    return (value === true || value === false);
};

let getType = function (data) {
    if (isObject(data)) {
        return 'object';
    } else if (isArray(data)) {
        return 'array';
    } else if (isNull(data)) {
        return null;
    } else if (isBoolean(data)) {
        return 'boolean';
    } else if (isString(data)) {
        return 'string';
    } else if (isNumber(data)) {
        return 'number';
    }
};

export {isObject, isArray, isNumber, isString, isBoolean, getType}
