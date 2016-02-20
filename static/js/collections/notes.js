/**
 * @fileoverview Describes a collection of notes
 */
define([
  'backbone',
  'models/note'
],
function(Backbone, NoteModel) {
  return Backbone.Collection.extend({
    url: '/api/notes/list',
    model: NoteModel
  });
});

