'use strict';

import ChangeEmitter from '../utils/ChangeEmitter.js';

import AppDispatcher from '../dispatcher/AppDispatcher';

import SearchConstants from '../constants/SearchConstants.js';


var _results = [];

function setResults(results) {
  if (results === null) {
    results = [];
  }
  _results = results;
}

function setInput(input) {
  localStorage.setItem("searchInput", input);
}

function input() {
  var s = localStorage.getItem("searchInput");
  if (s === null) {
    return "";
  }
  return s;
}


class SearchResultStore extends ChangeEmitter {
  getResults() {
    return _results;
  }

  getInput() {
    return input();
  }
}

var _searchResultStore = new SearchResultStore();

_searchResultStore.dispatchToken = AppDispatcher.register(function(payload) {
  var action = payload.action;
  var source = payload.source;

  if (source === 'SERVER_ACTION') {
    switch (action.actionType) {
      case SearchConstants.SEARCH:
        setResults(action.data);
        _searchResultStore.emitChange();
        break;
    }
  }

  if (source === 'VIEW_ACTION') {
    switch (action.actionType) {
      case SearchConstants.SEARCH:
        setInput(action.input);
        _searchResultStore.emitChange();
        break;
    }
  }

  return true;
});

export default _searchResultStore;
