// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

library IterableMapping {
    struct Map {
        bytes[] keys;
        mapping(bytes => CoordinatorInfo) values;
        mapping(bytes => uint) indexOf;
        mapping(bytes => bool) inserted;
    }

    struct CoordinatorInfo {
        address ownerAddress;
        address[] functionAddresses;
    }

    function get(Map storage map, bytes memory key) public view returns (CoordinatorInfo storage) {
        return map.values[key];
    }

    function getKeyAtIndex(Map storage map, uint index) public view returns (bytes memory) {
        return map.keys[index];
    }

    function size(Map storage map) public view returns (uint) {
        return map.keys.length;
    }

    function set(Map storage map, bytes memory key, CoordinatorInfo memory val) public {
        if (map.inserted[key]) {
            map.values[key] = val;
        } else {
            map.inserted[key] = true;
            map.values[key] = val;
            map.indexOf[key] = map.keys.length;
            map.keys.push(key);
        }
    }

    function remove(Map storage map, bytes memory key) public {
        if (!map.inserted[key]) {
            return;
        }

        delete map.inserted[key];
        delete map.values[key];

        uint index = map.indexOf[key];
        bytes memory lastKey = map.keys[map.keys.length - 1];

        map.indexOf[lastKey] = index;
        delete map.indexOf[key];

        map.keys[index] = lastKey;
        map.keys.pop();
    }
}