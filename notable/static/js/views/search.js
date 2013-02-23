/**
 * @fileoverview Describes a search modal view
 */
define([
  'text!templates/search.html',
  'backbone',
  'underscore',
  'lib/jquery.hotkeys'
],
function(searchTemplate) {
  return Backbone.View.extend({

    events: {
      'keyup .search-query': 'search',
      'submit form': 'noop'
    },

    hide: function() {
      this.trigger('search', null);
      this.$('input').animate({
        right: '-450px'
      }, 'fast', function() {
        $(this).val('');
      });
    },

    initialize: function() {
      this._query = null;
      $(document).bind('click', _.bind(this.hide, this));
      $(document).bind('keydown', '/', _.bind(this.show, this));
    },

    next: function() {
      this.trigger('next');
      return false;
    },

    previous: function() {
      this.trigger('previous');
      return false;
    },

    open: function() {
      this.trigger('open');
      this.hide();
      this.$('input').blur(); // (un)focus search
      return false;
    },

    noop: function() {
      return false;
    },

    render: function(collection) {
      this.$el.html(_.template(searchTemplate));
      this.$('.search-query').bind('keydown', 'ctrl+j', _.bind(this.next, this));
      this.$('.search-query').bind('keydown', 'ctrl+k', _.bind(this.previous, this));
      this.$('.search-query').bind('keydown', 'return', _.bind(this.open, this));
      return this;
    },

    show: function() {
      $('.nav-tabs a:first').tab('show');
      this.$('input').animate({
        right: '10px'
      }, 'fast', function() {
        $(this).focus();
      });
    },

    search: function() {
      var query = this.$('input').val().trim().toLowerCase()
      if (query != this._query) {
        this.trigger('search', query);
        this._query = query;
      }
    }

  });
});
