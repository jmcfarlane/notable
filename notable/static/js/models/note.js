/**
 * @fileoverview Describes a note
 */
define([
  'backbone'
],
function(Backbone) {
  return Backbone.Model.extend({

    url: function(atts, options) {
      return '/api/note/' + (this.get('uid') || 'create');
    },

    fetchContent: function(password) {
      $.ajax({
        url: '/api/note/content/' + this.get('uid'),
        type: 'POST',
        data: {
          password: password
        },
        success: _.bind(function(response, textStatus, xhr) {
          this.set({
            content: response,
            password: password
          });
        }, this),
        error: _.bind(function(xhr, response) {
          this.trigger('decryption:error', xhr.responseText);
        }, this)
      });
    }

  });
});

