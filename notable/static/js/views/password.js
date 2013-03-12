/**
 * @fileoverview Describes a password modal
 */
define([
  'text!templates/password.html',
  'backbone',
  'underscore'
],
function(passwordModalTemplate) {
  return Backbone.View.extend({
    events: {
      'click .modal-footer .btn:first': 'hide',
      'click .modal-footer .btn-primary': 'submit',
      'submit form': 'submit',
      'shown': 'setFocus'
    },

    getModal: function() {
      return this.$('div').first();
    },

    getPassword: function() {
      return this._modal.find('input');
    },

    hide: function() {
      return this._modal.modal('hide');
    },

    render: function(collection) {
      this.$el.html(_.template(passwordModalTemplate));
      this._modal = this.getModal();
      return this;
    },

    renderError: function(msg) {
      return this.$('.error').html(msg).show();
    },

    reset: function() {
      this.$('.error').html('').hide()
    },

    setFocus: function() {
      this.$('input').focus();
    },

    show: function(callback) {
      this.reset();
      this.callback = callback;
      this._modal.modal({
        backdrop: 'static',
        keyboard: true
      });
    },

    submit: function() {
      this.callback();
      this.$('input').val('');
      return false;
    }

  });
});
