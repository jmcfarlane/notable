/**
 * @fileoverview Describes a note
 */
define([
  'backbone'
],
function(Backbone) {
  return Backbone.Model.extend({
    getContent: function(password, callback) {
      $.ajax({
        url: '/api/note/content/' + this.get('uid'),
        type: 'POST',
        data: {password: password},
        success: function(response, textStatus, xhr) {
          callback(response);
        }
      }, this);
    }

  });
});

