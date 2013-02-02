/**
 * @fileoverview Describes a collection of notes
 */
define([
  'backbone'
],
function(Backbone) {
  return Backbone.Collection.extend({
    url: '/api/notes/list'
    //parse: function(response) {
    //  return _.map(response.rows, function(row) {
    //    var note = {};
    //    _.each(response.cols, function(element, index, list) {
    //      console.log('element:', element, index);
    //    });
    //    return new Backbone.Model(
    //      _.extend(row, response.cols
    //    ));
    //  });
    //}
  });
});

