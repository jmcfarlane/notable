/**
 * @fileoverview Describes an application usage modal
 */
define([
  'text!templates/usage.html',
  'text!templates/version.html',
  'backbone',
  'underscore'
],
function(usageModalTemplate, versionTemplate) {
  return Backbone.View.extend({

    initialize: function() {
      var Model = Backbone.Model.extend({ url: '/api/version' });
      this.model = new Model();
      Mousetrap.bind('?', _.bind(this.show, this));
    },

    events: {
      'click .modal-footer .btn:first': 'hide'
    },

    getModal: function() {
      return this.$('div').first();
    },

    hide: function() {
      return this._modal.modal('hide');
    },

    render: function() {
      this.$el.html(_.template(usageModalTemplate));
      this._modal = this.getModal();
      return this;
    },


    show: function(callback) {
      this.model.fetch({
        success: _.bind(function(model, response, options){
          this.$el.find("#version-info").html(_.template(versionTemplate, {
            model: this.model.toJSON()
          }));
        }, this),
      });
      this._modal.modal({
        backdrop: 'static',
        keyboard: false
      });
    },

  });
});
