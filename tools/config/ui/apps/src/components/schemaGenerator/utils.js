'use strict';

let isObject = function(value) {
  return (null !== value && typeof value === typeof {} && !isArray(value));
};

let isNumber = function(value) {
  return !isArray( value ) && (value - parseFloat( value ) + 1) >= 0;
};

let isArray = function(value) {
  return (value instanceof Array);
};

let isString = function(value) {
  return (typeof value === typeof '');
};

let isNull = function(value) {
  return (null === value);
};

let isBoolean = function(value) {
  return (value === true || value === false);
};

let toObject = function(arr) {
  var rv = {};
  for (var i = 0; i < arr.length; ++i)
    rv[i] = arr[i];
  return rv;
};

let oneIsNull = function(v1, v2) {
  return ((v1 === null && v2 !== null) || (v1 !== null && v2 === null));
};

let isUndefined = function(val) {
  return (null === val || typeof val === 'undefined');
};

let isFunction = function(fn) {
  return (typeof fn === 'function');
};

let isEqual = function(v1, v2) {
  if (typeof v1 !== typeof v2 || oneIsNull(v1, v2)) {
    return false;
  }

  if (typeof v1 === typeof "" || typeof v1 === typeof 0) {
    return v1 === v2;
  }

  var _isEqual = true;

  if (typeof v1 === typeof {}) {
    var compare = function(value1, value2) {
      for (var i in value1) {
        if (!value2.hasOwnProperty(i)) {
          _isEqual = false;
          break;
        }

        if (exports.isObject(value1[i])) {
          compare(value1[i], value2[i]);
        } else if (typeof value1[i] === typeof "") {
          if (value1[i] !== value2[i]) {
            _isEqual = false;
            break;
          }
        }
      }
    }

    compare(v1, v2);
  }

  return _isEqual;
};

let getType = function(data) {
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

export {isObject,isNumber,isArray,isString,isNull,isBoolean,toObject,oneIsNull,isUndefined,isFunction,isEqual,getType }
