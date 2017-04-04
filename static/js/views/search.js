/**
 * @fileoverview Describes a search modal view
 */
define([
  'text!templates/search.html',
  'lib/mousetrap.min',
  'backbone',
  'underscore'
],
function(searchTemplate, Mousetrap) {
  return Backbone.View.extend({

    events: {
      'keyup #search': 'search',
      'submit form': 'noop',
      'focus input': 'activatePrimaryTab'
    },

    initialize: function() {
      this._query = null;
      Mousetrap.bind('/', _.bind(this.show, this));
    },

    activatePrimaryTab: function() {
      $('.nav-tabs a:first').tab('show');
    },

    focus: function() {
      this.$('input').focus();
    },

    isMidSearch: function() {
      return this.$('input').val().length > 0;
    },

    next: function() {
      this.trigger('next');
      return false;
    },

    previous: function() {
      this.trigger('previous');
      return false;
    },

    noop: function() {
      return false;
    },

    render: function(collection) {
      this.$el.html(_.template(searchTemplate));
      Mousetrap.bind('ctrl+j', _.bind(this.next, this));
      Mousetrap.bind('ctrl+k', _.bind(this.previous, this));
      return this;
    },

    show: function(e) {
      if (typeof document.activeElement.form != 'undefined') {
        var el = $(e.target);
        el.val(el.val() + e.key);
        return false;
      }
      this.focus();
      return false
    },

    search: function() {
      var query = this.$('input').val().trim().toLowerCase()
      if (query != this._query) {
        $.ajax({
          url: '/api/notes/search',
          type: 'get',
          data: {
            q: query
          },
          error: function(xhr) {
            console.log('Error:', xhr);
            this.trigger('search', {
              query: query,
              uids: []
            });
            this._query = query;
          }.bind(this),
          success: function(resp) {
            this.trigger('search', {
              query: query,
              uids: resp
            });
            this._query = query;
          }.bind(this)
        });
      }
    }

  });
});
