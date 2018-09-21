/**
 * @fileoverview Describes a confirmation modal
 */
define([
  'text!templates/confirm.html',
  'backbone',
  'underscore'
],
function(confirmationModalTemplate) {
  return Backbone.View.extend({
    events: {
      'click .modal-footer .btn:first': 'hide',
      'click .modal-footer .btn-danger': 'proceed',
    },

    getModal: function() {
      return this.$('div').first();
    },

    hide: function() {
      return this._modal.modal('hide');
    },

    render: function(action, msg) {
      this.$el.html(_.template(confirmationModalTemplate));
      this._modal = this.getModal();
      $(this.$el).find('h3').html(msg)
      $(this.$el).find('.btn-danger').html(action)
      return this;
    },

    show: function(callback) {
      this.callback = callback;
      this._modal.modal({
        backdrop: 'static'
      });
    },

    proceed: function() {
      this.callback();
      return false;
    }

  });
});
