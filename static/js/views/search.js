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
      'keyup .search-query': 'search',
      'submit form': 'noop'
    },

    hide: function() {
      this.trigger('search', null);
      this.$('input').animate({
        right: '-450px'
      }, 'fast', function() {
        $(this).val('').blur();
      });
    },

    initialize: function() {
      this._query = null;
      $(document).bind('click', _.bind(this.hide, this));
      Mousetrap.bind('/', _.bind(this.show, this));
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
      Mousetrap.bind('ctrl+j', _.bind(this.next, this));
      Mousetrap.bind('ctrl+k', _.bind(this.previous, this));
      Mousetrap.bind('esc', _.bind(this.hide, this));
      return this;
    },

    show: function() {
      $('.nav-tabs a:first').tab('show');
      this.$('input').animate({
        right: '0px'
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
