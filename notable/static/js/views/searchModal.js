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
      $(document).bind('click', _.bind(this.hide, this));
      $(document).bind('keydown', '/', _.bind(this.show, this));
    },

    noop: function() {
      return false;
    },

    render: function(collection) {
      this.$el.html(_.template(searchTemplate));
      return this;
    },

    show: function() {
      this.$('input').animate({
        right: '10px'
      }, 'fast', function() {
        $(this).focus();
      });
    },

    search: function() {
      this.trigger('search', this.$('input').val().trim().toLowerCase());
    }

  });
});
