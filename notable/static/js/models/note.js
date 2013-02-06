/**
 * @fileoverview Describes a note
 */
define([
  'backbone'
],
function(Backbone) {
  return Backbone.Model.extend({
    fetchContent: function(password) {
      $.ajax({
        url: '/api/note/content/' + this.get('uid'),
        type: 'POST',
        data: {
          password: password
        },
        success: _.bind(function(response, textStatus, xhr) {
          this.set('content', response);
        }, this),
        error: _.bind(function(xhr, response) {
          this.trigger('decryption:error', xhr.responseText);
        }, this)
      });
    }

  });
});

